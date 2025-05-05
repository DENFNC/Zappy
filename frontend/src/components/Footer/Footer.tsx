import SocialLinks from '../SocialLinks/SocialLinks';

export default function Footer() {
  return (
    <footer className="bg-white text-gray-700 border-t border-gray-200">
      <div className="max-w-7xl mx-auto px-4 py-10 sm:px-6 lg:px-8">
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-8">
          <div>
            <h3 className="text-sm font-semibold mb-3">Компания</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <a href="#" className="hover:text-blue-600">
                  О нас
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-blue-600">
                  Блог
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-blue-600">
                  Вакансии
                </a>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-sm font-semibold mb-3">Поддержка</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <a href="#" className="hover:text-blue-600">
                  Связаться с нами
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-blue-600">
                  FAQ
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-blue-600">
                  Помощь
                </a>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-sm font-semibold mb-3">Юридическая информация</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <a href="#" className="hover:text-blue-600">
                  Политика конфиденциальности
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-blue-600">
                  Условия использования
                </a>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-sm font-semibold mb-3">Ресурсы</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <a href="#" className="hover:text-blue-600">
                  Руководства
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-blue-600">
                  Партнёры
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-blue-600">
                  API
                </a>
              </li>
            </ul>
          </div>
          <div>
            <SocialLinks />
          </div>
        </div>
        <div className="mt-10 text-sm text-center text-gray-500">
          © {new Date().getFullYear()} Ваша Компания. Все права защищены.
        </div>
      </div>
    </footer>
  );
}
