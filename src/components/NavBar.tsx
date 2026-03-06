import { useState } from "react";

interface Props {
  onSelectItem: (item: string) => void;
}

const NavBar = ({ onSelectItem }: Props) => {
  const entries = ["Jobs", "Workers"];
  const [selectedIndex, setSelectedIndex] = useState(0);

  return (
    <>
      <ul className="list-group list-group-horizontal">
        {entries.map((item, index) => (
          <li
            key={index}
            className={
              selectedIndex === index
                ? "list-group-item active"
                : "list-group-item"
            }
            onClick={() => {
              setSelectedIndex(index);
              onSelectItem(item);
            }}
          >
            {item}
          </li>
        ))}
      </ul>
    </>
  );
};

export default NavBar;
