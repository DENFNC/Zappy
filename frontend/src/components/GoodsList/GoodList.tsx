import { CardGoodType } from '@/types/card-goods.type';
import CardGood from '../CardGood/CardGood';

interface Props {
  title: string;
  cards: CardGoodType[];
}

export default function GoodList({ title, cards }: Props) {
  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold mb-6">{title}</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {cards.map((card) => (
          <CardGood
            image={card.image}
            description={card.description}
            price={card.price}
            name={card.name}
            key={card.id}
            id={card.id}
          />
        ))}
      </div>
    </div>
  );
}
