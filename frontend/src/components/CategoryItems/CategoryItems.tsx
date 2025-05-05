import { useState } from 'react';
import { useResponsiveLimit } from '@/hooks/useResponsiveLimit';

export interface IListItems {
  items: string[];
  limit?: number;
}

export default function CategoryItems({ items, limit = 10 }: IListItems) {
  const responsiveLimit = useResponsiveLimit(limit, {
    640: 3, // for screens <= 640px
    768: 4, // for screens <= 768px
    1024: 5, // for screens <= 1024px
  });

  const [isDropdownOpen, setDropdownOpen] = useState(false);

  const visibleItems = items.slice(0, responsiveLimit);
  const hiddenItems = items.slice(responsiveLimit);

  return (
    <div className="w-full px-4 py-2">
      <ul className="flex flex-wrap justify-center gap-3 sm:gap-4 md:gap-6">
        {visibleItems.map((item, index) => (
          <li
            key={index}
            className="text-sm text-gray-500 hover:text-yellow-700 cursor-pointer font-extrabold"
          >
            {item}
          </li>
        ))}

        {hiddenItems.length > 0 && (
          <li className="relative items-center">
            <button
              onClick={() => setDropdownOpen(!isDropdownOpen)}
              className="text-sm text-gray-500 hover:text-yellow-700 cursor-pointer font-bold"
            >
              Ещё ▾
            </button>
            {isDropdownOpen && (
              <ul className="absolute z-10 mt-2 w-40 bg-white border border-gray-200 rounded shadow-lg">
                {hiddenItems.map((item, index) => (
                  <li
                    key={index}
                    className="px-4 py-2 text-sm text-gray-700 hover:bg-yellow-100 hover:text-yellow-700 cursor-pointer"
                  >
                    {item}
                  </li>
                ))}
              </ul>
            )}
          </li>
        )}
      </ul>
    </div>
  );
}
