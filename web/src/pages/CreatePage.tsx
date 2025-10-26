import { useNavigate } from 'react-router-dom'
import { ProductForm } from '../components/product-form'
import { productClient } from '../grpc'

export default function CreatePage() {
  const navigate = useNavigate()

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-md rounded border p-6 shadow-sm bg-white">
        <h2 className="mb-4 text-xl font-bold">Create Product</h2>
        <ProductForm
          onSubmit={async (p) => {
            try {
              await productClient.createProduct({
                name: p.name,
                description: p.description,
                price: p.price,
                cover: p.cover,
              })
              navigate('/')
            } catch (e) {
              console.error(e)
            }
          }}
        />
      </div>
    </div>
  )
}
