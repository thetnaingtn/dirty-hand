import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { ProductDetail } from '../components/product-detail'
import type { Product } from '../types/proto/api/v1/product'
import { productClient } from '../grpc'

export default function DetailPage() {
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
    <div className="p-4">
      <ProductDetail
        product={product}
        onEdit={() => navigate(`/products/${product.id.toString()}/edit`)}
        onBack={() => navigate('/')}
      />
    </div>
  )
}
