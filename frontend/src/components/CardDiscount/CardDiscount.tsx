interface Props {
  title: string;
  description: string;
}

export default function Card(props: Props) {
  return (
    <div className="max-w-sm rounded-2xl overflow-hidden shadow-lg border border-gray-200 bg-white w-full">
      <div className="p-4">
        <h2 className="text-xl font-semibold text-gray-800 mb-2">{props.title}</h2>
        <p className="text-gray-600 text-sm">{props.description}</p>
      </div>
    </div>
  );
}
