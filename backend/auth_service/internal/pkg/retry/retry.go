package retry

import (
	"context"
	"errors"
	"time"
)

// BackoffPolicy определяет, как вычислять задержку перед next retry.
// Принимает номер попытки (начиная с 1) и возвращает длительность сна.
type BackoffPolicy func(attempt int) time.Duration

// ShouldRetry определяет, стоит ли пробовать снова при данной ошибке.
// Если возвращает true — следующий retry, иначе — сразу выходим с ошибкой.
type ShouldRetry func(err error) bool

// Option задаёт какую-то опцию для Retry (functional option).
type Option func(cfg *retryConfig)

// retryConfig хранит в себе все настройки, выбранные через Option.
type retryConfig struct {
	maxRetries  int
	backoff     BackoffPolicy
	shouldRetry ShouldRetry
}

// DefaultMaxRetries — если опцию WithMaxRetries не указывать, возьмётся 3 retry‐attempts.
const DefaultMaxRetries = 3

// DefaultBackoff сразу возвращает 0, то есть нет задержки между попытками.
func DefaultBackoff(_attempt int) time.Duration {
	return 0
}

// DefaultShouldRetry всегда возвращает true: переходим к следующей попытке при любой ошибке.
func DefaultShouldRetry(_err error) bool {
	return true
}

// WithMaxRetries задаёт максимальное число дополнительных попыток (после первой).
// Если передать n < 0, будет возвращена ошибка при вызове Retry.
func WithMaxRetries(n int) Option {
	return func(cfg *retryConfig) {
		cfg.maxRetries = n
	}
}

// WithInterval задаёт фиксированную задержку между попытками (linear backoff).
// За кулисами это всего лишь backoff-policy, которая не зависит от attempt.
func WithInterval(d time.Duration) Option {
	return func(cfg *retryConfig) {
		cfg.backoff = func(_attempt int) time.Duration {
			return d
		}
	}
}

// WithBackoff задаёт кастомную стратегию экспоненциального (или произвольного) backoff.
// Приоритетнее, чем WithInterval, если обе указаны: просто перезапишет backoff.
func WithBackoff(policy BackoffPolicy) Option {
	return func(cfg *retryConfig) {
		if policy != nil {
			cfg.backoff = policy
		}
	}
}

// WithShouldRetry задаёт, какие ошибки позволено повторять.
// Если функция возвращает false для err, retry прекращается сразу.
func WithShouldRetry(fn ShouldRetry) Option {
	return func(cfg *retryConfig) {
		if fn != nil {
			cfg.shouldRetry = fn
		}
	}
}

// Retry выполняет doFunc не более maxRetries+1 раз:
//   - ctx: контекст, отмена которого прерывает все последующие попытки (возврат ctx.Err()).
//   - doFunc: func(ctx context.Context) (R, error).
//   - opts: список опций (WithMaxRetries, WithInterval, WithBackoff, WithShouldRetry).
//
// Порядок применения опций:
//  1. maxRetries по умолчанию = DefaultMaxRetries (3).
//  2. backoff по умолчанию = DefaultBackoff (то есть сразу, без Sleep).
//  3. shouldRetry по умолчанию = DefaultShouldRetry (то есть retry на любой err).
//
// Возвращает либо «успешный R», либо последний error (или ctx.Err()).
func Retry[R any](
	ctx context.Context,
	doFunc func(ctx context.Context) (R, error),
	opts ...Option,
) (R, error) {
	cfg := &retryConfig{
		maxRetries:  DefaultMaxRetries,
		backoff:     DefaultBackoff,
		shouldRetry: DefaultShouldRetry,
	}

	for _, o := range opts {
		o(cfg)
	}

	var zero R
	if cfg.maxRetries < 0 {
		return zero, errors.New("retry: MaxRetries must be >= 0")
	}

	var lastErr error
	for attempt := 0; attempt <= cfg.maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		default:
		}

		res, err := doFunc(ctx)
		if err == nil {
			return res, nil
		}

		if cfg.shouldRetry != nil && !cfg.shouldRetry(err) {
			return zero, err
		}

		lastErr = err
		if attempt < cfg.maxRetries {
			sleepDur := cfg.backoff(attempt + 1)
			time.Sleep(sleepDur)
		}
	}

	return zero, lastErr
}
