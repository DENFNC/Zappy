export default function Header() {
  return (
    <header className="bg-white">
      <div className="mx-auto px-4 p-2">
        <div className="flex justify-between justify-center h-16 items-center ml">
          <h2 className="text-purple-600 text-xl">Zappy</h2>

          <nav className="hidden md:flex space-x-6">
            <a href="#" className="text-gray-700 hover:text-purple-600">
              Войти
            </a>
            <a href="#" className="text-gray-700 hover:text-purple-600">
              Избранное
            </a>
            <a href="#" className="text-gray-700 hover:text-purple-600">
              Корзина
            </a>
          </nav>
        </div>
        <div className="bottom-header-wrapper">
          <div className="flex justify-between justify-center h-16 items-center">Bottom header</div>
        </div>
      </div>
    </header>
  );
}
