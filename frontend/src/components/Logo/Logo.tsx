import Image from 'next/image';

interface Props {
  width: number;
  height?: number;
}

export default function Logo(props: Props) {
  return (
    <Image
      className="mr-2"
      height={props.height ?? 20}
      width={props.width}
      alt={'Zappy'}
      src={'/logo.png'}
    ></Image>
  );
}
