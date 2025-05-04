import { CardGoodType } from '@/types/card-goods.type';
import Image from 'next/image';

export default function CardGood(props: CardGoodType) {
  return (
    <div
      className="max-w-sm bg-white rounded-2xl shadow-md overflow-hidden hover:shadow-xl transition-shadow duration-300"
      key={props.id}
    >
      <Image
        className="w-full h-48 object-cover"
        src={props.image}
        alt={props.name}
        width={300} // Specify width
        height={200}
      />
      <div className="p-4">
        <h2 className="text-xl font-semibold text-gray-800">{props.name}</h2>
        <p className="text-gray-600 mt-2">{props.description}</p>
        <div className="mt-4 flex items-center justify-between">
          <span className="text-lg font-bold text-yellow-800">${props.price}</span>
          <button className="px-4 py-2 bg-gray-400 text-white text-sm font-medium rounded hover:bg-yellow-700">
            Add to Cart
          </button>
        </div>
      </div>
    </div>
  );
}
