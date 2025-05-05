import CategoryItems from '../CategoryItems/CategoryItems';
import LogoSvg from '../Logo/Logo';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faUser,
  faHeart,
  faShoppingCart,
  faBars,
  faSearch,
} from '@fortawesome/free-solid-svg-icons';

export default function Header() {
  return (
    <>
      <header className="bg-white top-0 z-50">
        <div className="h-16 flex items-center justify-center px-2 sm:px-4">
          <nav className="flex w-full items-center space-x-2 sm:space-x-4 overflow-x-auto">
            <LogoSvg width={40} />

            <button
              className="bg-gray-300 hover:bg-yellow-700 text-white p-2 rounded h-10 w-10 flex items-center justify-center"
              aria-label="Меню"
            >
              <FontAwesomeIcon icon={faBars} className="h-5 w-5" />
            </button>

            {/* Search bar: takes remaining space, shrinks on small screens */}
            <div className="relative flex items-center flex-grow max-w-full border border-gray-200 rounded-lg bg-white">
              <input
                type="text"
                className="h-10 w-full px-10 pr-4 rounded-lg focus:outline-none"
                placeholder="Искать товары и категории"
              />
              <FontAwesomeIcon
                icon={faSearch}
                className="absolute left-3 text-gray-400 pointer-events-none"
              />
              <button className="h-10 px-4 bg-gray-300 hover:bg-yellow-700 text-white flex items-center justify-center rounded-r-lg">
                <FontAwesomeIcon icon={faSearch} className="h-4 w-4" />
              </button>
            </div>

            {/* Buttons shrink but stay in row */}
            <div className="flex shrink-0 space-x-2">
              <button
                className="bg-gray-300 hover:bg-yellow-700 text-white p-2 rounded h-10 w-10"
                aria-label="Войти"
              >
                <FontAwesomeIcon icon={faUser} className="h-5 w-5" />
              </button>
              <button
                className="bg-gray-300 hover:bg-yellow-700 text-white p-2 rounded h-10 w-10"
                aria-label="Избранное"
              >
                <FontAwesomeIcon icon={faHeart} className="h-5 w-5" />
              </button>
              <button
                className="bg-gray-300 hover:bg-yellow-700 text-white p-2 rounded h-10 w-10"
                aria-label="Корзина"
              >
                <FontAwesomeIcon icon={faShoppingCart} className="h-5 w-5" />
              </button>
            </div>
          </nav>
        </div>

        {/* Keep CategoryItems visible only on larger screens if needed */}
        <CategoryItems
          items={[
            'Электроника',
            'Бытовая техника',
            'Одежда',
            'Красота и уход',
            'Здоровье',
            'Товары для дома',
            'Спортивные товары',
            'Автотовары',
            'Игрушки и игры',
            'Книги',
            'Товары для детей',
            'Продукты питания',
            'Мебель',
            'Обувь',
            'Украшения',
            'Товары для дачи',
          ]}
        />
      </header>
    </>
  );
}
