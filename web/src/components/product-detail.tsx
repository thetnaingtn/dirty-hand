import { Button } from './ui/button'
import type { Product } from '../types/proto/api/v1/product'

interface Props {
  product: Product
  onEdit: () => void
  onBack: () => void
}

export function ProductDetail({ product, onEdit, onBack }: Props) {
  return (
    <div className="space-y-4">
      <h2 className="text-xl font-bold">{product.name}</h2>
      <p>Price: {product.price}</p>
      <div className="space-x-2">
        <Button onClick={onEdit}>Edit</Button>
        <Button onClick={onBack} className="bg-gray-200 text-black hover:bg-gray-300">Back</Button>
      </div>
    </div>
  )
}
