import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { ProductForm } from '../components/product-form'
import type { Product } from '../types/proto/api/v1/product'
import { productClient } from '../grpc'

export default function UpdatePage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [product, setProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      try {
        const res = await productClient.listProducts({})
        const found = res.products.find((p) => p.id.toString() === id)
        setProduct(found ?? null)
      } catch (e) {
        console.error(e)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [id])

  if (loading) return <div className="p-4">Loading...</div>
  if (!product) return (
    <div className="p-4 space-y-2">
      <div>Product not found.</div>
      <button className="underline" onClick={() => navigate('/')}>Back to list</button>
    </div>
  )

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-md rounded border p-6 shadow-sm bg-white">
        <h2 className="mb-4 text-xl font-bold">Update Product</h2>
        <ProductForm
          initial={{ id: product.id, name: product.name, description: product.description, price: product.price, cover: product.cover }}
          onSubmit={async (p) => {
            try {
              await productClient.updateProduct({
                id: p.id ?? product.id,
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
