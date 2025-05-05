import { useState } from 'react';

export interface IListItems {
  items: string[];
  limit?: number;
}

export default function CategoryItems({ items, limit = 6 }: IListItems) {
  const [isDropdownOpen, setDropdownOpen] = useState(false);

  const visibleItems = items.slice(0, limit);
  const hiddenItems = items.slice(limit);

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
          <li className="relative">
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
