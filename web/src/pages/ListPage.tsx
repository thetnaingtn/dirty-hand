import { useNavigate } from 'react-router-dom'
import { ProductList } from '../components/product-list'

export default function ListPage() {
  const navigate = useNavigate()

  return (
    <div className="p-4">
      <ProductList
        onSelect={(p) => navigate(`/products/${p.id.toString()}`)}
        onCreate={() => navigate('/products/new')}
        onUpdate={(p) => navigate(`/products/${p.id.toString()}/edit`)}
      />
    </div>
  )
}
