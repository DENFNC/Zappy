import CategoryItems from '../CategoryItems/CategoryItems';
import LogoSvg from '../Logo/Logo';

export default function Header() {
  return (
    <>
      <header className="bg-white sticky top-0 z-50">
        <div className="h-16 flex items-center justify-center">
          <nav className="flex w-full items-center space-x-6">
            <LogoSvg width={40} />

            <button className="bg-gray-300 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded h-10">
              Каталог
            </button>

            <div className="flex border-1 border-gray-200 rounded w-full">
              <input
                type="text"
                className="bg-white h-10 w-full min-w-80 px-2 rounded-lg focus:outline-none hover:cursor-pointer rounded"
                name=""
                placeholder="Искать товары и категории"
              />
            </div>

            <div className="flex m-auto space-x-4">
              <button className="flex-1 bg-gray-300 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded h-10">
                Войти
              </button>

              <button className="flex-1 bg-gray-300 hover:hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded h-10">
                Избранное
              </button>

              <button className="flex-1 bg-gray-300 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded h-10">
                Корзина
              </button>
            </div>
          </nav>
        </div>

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
          limit={11}
        />
      </header>
    </>
  );
}
