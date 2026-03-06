interface Props {
  items: Array<any>;
}

const ListGroup = ({ items }: Props) => {
  console.log(items);
  return (
    <>
      <ul className="list-group">
        {items.map((item) => (
          <li key={item.ID} className="list-group-item">
            {item.ID} | {item.Status}
          </li>
        ))}
      </ul>
    </>
  );
};

export default ListGroup;
