import { useState } from "react";

interface Props {
  onSelectItem: (item: string) => void;
}

const NavBar = ({ onSelectItem }: Props) => {
  const entries = ["Jobs", "Workers"];
  const [selectedIndex, setSelectedIndex] = useState(0);

  return (
    <div className="list-group list-group-horizontal">
      {entries.map((item, index) => (
        <a
          href="#"
          key={index}
          className={
            selectedIndex === index
              ? "list-group-item list-group-item-action active"
              : "list-group-item list-group-item-action"
          }
          onClick={() => {
            setSelectedIndex(index);
            onSelectItem(item);
          }}
        >
          {item}
        </a>
      ))}
    </div>
  );
};

export default NavBar;
