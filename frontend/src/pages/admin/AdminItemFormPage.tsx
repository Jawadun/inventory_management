import { ItemForm } from '../../components/common/ItemForm';

export function AdminItemFormPage({ isEdit = false }: { isEdit?: boolean }) {
  return <ItemForm isEdit={isEdit} />;
}