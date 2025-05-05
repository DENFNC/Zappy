'use client';
import CardDiscount from '@/components/CardDiscount/CardDiscount';
import CarouselDiscountItems from '@/components/CarouselDiscountItems/CarouselDiscountItems';
import Footer from '@/components/Footer/Footer';
import Header from '@/components/Header/Header';

import GoodList from '@/components/GoodsList/GoodList';
import { CardGoodType } from '@/types/card-goods.type';

const goods: CardGoodType[] = [
  {
    id: '1',
    name: 'Wireless Headphones',
    description: 'High-quality sound with noise cancellation.',
    image:
      'https://images.unsplash.com/photo-1603791440384-56cd371ee9a7?crop=entropy&cs=tinysrgb&fit=max&ixid=MXwyMDg5OHwwfDF8c2VhY2h8MXx8aGVhZHBob25lfGVufDB8fHx8fDE2NTg4NDA0NjI&ixlib=rb-1.2.1&q=80&w=400',
    price: 89.99,
  },
  {
    id: '2',
    name: 'Smart Watch',
    description: 'Stay connected with fitness tracking.',
    image:
      'https://images.unsplash.com/photo-1571944692545-e34e2cf80065?crop=entropy&cs=tinysrgb&fit=max&ixid=MXwyMDg5OHwwfDF8c2VhY2h8MnwxfGZpdG5lc3MlMkZyZWQgdGVsZXZpc2lvbnxlbnwwfDF8fHx8fDE2NjI0NzYzMzc&ixlib=rb-1.2.1&q=80&w=400',
    price: 129.99,
  },
  {
    id: '3',
    name: 'Bluetooth Speaker',
    description: 'Portable speaker with powerful bass.',
    image:
      'https://images.unsplash.com/photo-1591101215045-c73e26bff53f?crop=entropy&cs=tinysrgb&fit=max&ixid=MXwyMDg5OHwwfDF8c2VhY2h8Mnx8YnJlYWtmYXJlZCBzcGVha2VydxlhbGwxfDB8fHx8fDE2NjI0NzYzMzc&ixlib=rb-1.2.1&q=80&w=400',
    price: 49.99,
  },
  {
    id: '4',
    name: 'Laptop Backpack',
    description: 'Durable backpack for your laptop and accessories.',
    image:
      'https://images.unsplash.com/photo-1573185599930-44128306e6b7?crop=entropy&cs=tinysrgb&fit=max&ixid=MXwyMDg5OHwwfDF8c2VhY2h8Mnx8Y2FwdHVyZSBpbnNlY3R8ZW58MHx8fHx8fDE2NjI0NzY2NDA&ixlib=rb-1.2.1&q=80&w=400',
    price: 59.99,
  },
];

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-grow">
        <CarouselDiscountItems />
        <div className="flex mt-4 space-x-3 text-center">
          <CardDiscount
            title={'Детский мир'}
            description="Скидки на игрушки, одежду и товары для детей."
          />
          <CardDiscount
            title={'Выгода'}
            description="Товары по лучшим ценам для экономных покупателей."
          />
          <CardDiscount
            title={'Модный базар'}
            description="Актуальные тренды и модные новинки по привлекательным ценам."
          />
          <CardDiscount
            title={'Товары недели'}
            description="Лучшие предложения и акции этой недели."
          />
        </div>

        <div className="flex">
          <GoodList title="Модная распродажа" cards={goods}></GoodList>
        </div>

        <div className="flex">
          <GoodList title="Турецкий бренд посуды" cards={goods}></GoodList>
        </div>

        <div className="flex">
          <GoodList title="Ваша красота" cards={goods}></GoodList>
        </div>

        <div className="flex">
          <GoodList title="Холодильники" cards={goods}></GoodList>
        </div>

        <div className="flex">
          <GoodList title="Все для здоровья" cards={goods}></GoodList>
        </div>

        <div className="flex">
          <GoodList title="Бренды малышей" cards={goods}></GoodList>
        </div>

        <div className="flex">
          <GoodList title="Скидки" cards={goods}></GoodList>
        </div>
      </main>
      <Footer />
    </div>
  );
}
