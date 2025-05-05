export default function SocialLinks() {
  return (
    <div>
      <h3 className="text-sm font-semibold mb-3">Мы в соцсетях</h3>
      <div className="flex space-x-4">
        <a href="#" className="text-gray-500 hover:text-blue-600" aria-label="Facebook">
          <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M22 12c0-5.52-4.48-10-10-10S2 6.48 2 12c0 5 3.66 9.13 8.44 9.88v-6.99h-2.54V12h2.54V9.8c0-2.5 1.49-3.89 3.77-3.89 1.09 0 2.23.2 2.23.2v2.45h-1.26c-1.24 0-1.62.77-1.62 1.56V12h2.75l-.44 2.89h-2.31v6.99C18.34 21.13 22 17 22 12z" />
          </svg>
        </a>
        <a href="#" className="text-gray-500 hover:text-blue-600" aria-label="Instagram">
          <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M7.75 2h8.5C19.2 2 22 4.8 22 7.75v8.5C22 19.2 19.2 22 16.25 22h-8.5C4.8 22 2 19.2 2 16.25v-8.5C2 4.8 4.8 2 7.75 2zm0 2C5.68 4 4 5.68 4 7.75v8.5C4 18.32 5.68 20 7.75 20h8.5c2.07 0 3.75-1.68 3.75-3.75v-8.5C20 5.68 18.32 4 16.25 4h-8.5zM12 7c2.76 0 5 2.24 5 5s-2.24 5-5 5-5-2.24-5-5 2.24-5 5-5zm0 2a3 3 0 100 6 3 3 0 000-6zm4.75-.75a1 1 0 110 2 1 1 0 010-2z" />
          </svg>
        </a>
      </div>
    </div>
  );
}
