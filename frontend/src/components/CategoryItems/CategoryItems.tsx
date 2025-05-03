export interface IListItems {
  items: string[];
}

export default function CategoryItems({ items }: IListItems) {
  return (
    <div className="w-full h-16 flex items-center justify-center">
      <ul className="flex w-full space-x-6">
        {items.map((item, index) => (
          <li
            key={index}
            className="text-sm text-gray-500 hover:text-purple-500 hover:cursor-pointer font-extrabold"
          >
            {item}
          </li>
        ))}
      </ul>
    </div>
  );
}
