export interface IListItems {
  items: string[];
}

export default function CategoryItems({ items }: IListItems) {
  return (
    <div className="w-full px-4 py-2">
      <ul className="flex flex-wrap justify-center gap-3 sm:gap-4 md:gap-6">
        {items.map((item, index) => (
          <li
            key={index}
            className="text-sm text-gray-500 hover:text-yellow-700 cursor-pointer font-extrabold"
          >
            {item}
          </li>
        ))}
      </ul>
    </div>
  );
}
