import { useNavigate } from 'react-router-dom';
import { InventoryList } from '../../components/common/InventoryList';
import { Item } from '../../types';

export function InventoryPage() {
  const navigate = useNavigate();

  const handleItemSelect = (item: Item) => {
    navigate(`/inventory/${item.id}`);
  };

  return (
    <div>
      <InventoryList isAdmin={false} onItemSelect={handleItemSelect} />
    </div>
  );
}