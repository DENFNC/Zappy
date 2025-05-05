import { useState, useEffect } from 'react';

export function useResponsiveLimit(defaultLimit: number, breakpoints: { [width: number]: number }) {
  const [limit, setLimit] = useState(defaultLimit);

  useEffect(() => {
    const updateLimit = () => {
      const width = window.innerWidth;
      const sortedBreakpoints = Object.keys(breakpoints)
        .map(Number)
        .sort((a, b) => a - b);

      for (const bp of sortedBreakpoints) {
        if (width <= bp) {
          setLimit(breakpoints[bp]);
          return;
        }
      }

      setLimit(defaultLimit);
    };

    updateLimit();
    window.addEventListener('resize', updateLimit);
    return () => window.removeEventListener('resize', updateLimit);
  }, [breakpoints, defaultLimit]);

  return limit;
}
