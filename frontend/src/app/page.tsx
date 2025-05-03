'use client';
import Footer from '@/components/Footer/Footer';
import Header from '@/components/Header/Header';

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-grow">Main content</main>
      <Footer />
    </div>
  );
}
