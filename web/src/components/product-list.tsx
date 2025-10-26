import { useEffect, useState } from 'react'
import { Button } from './ui/button'
import { AiFillDelete, AiFillEdit } from 'react-icons/ai'
import type { Product } from '../types/proto/api/v1/product'
import { productClient } from '../grpc'

interface Props {
  onSelect: (product: Product) => void
  onCreate: () => void
  onUpdate: (product: Product) => void
}

export function ProductList({ onSelect, onCreate, onUpdate }: Props) {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(false)

  async function fetchProducts() {
    try {
      setLoading(true)
      const res = await productClient.listProducts({})
      setProducts(res.products)
    } catch (e) {
      console.error(e)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchProducts()
  }, [])

  const handleDelete = async (id: number) => {
    try {
      await productClient.deleteProduct({ id })
      setProducts((prev) => prev.filter((p) => p.id !== id))
    } catch (e) {
      console.error(e)
    }
  }

  const truncate = (s: string, n = 100) => (s.length > n ? s.slice(0, n) + '…' : s)

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-bold">Products</h2>
        <Button onClick={onCreate}>Create</Button>
      </div>
      <table className="w-full border">
        <thead>
          <tr className="bg-gray-100 text-left">
            <th className="p-2 border">Cover</th>
            <th className="p-2 border">Name</th>
            <th className="p-2 border">Description</th>
            <th className="p-2 border">Price</th>
            <th className="p-2 border w-[1%]">Actions</th>
          </tr>
        </thead>
        <tbody>
          {products.map((p) => (
            <tr key={p.id} className="hover:bg-gray-50">
              <td className="p-2 border" onClick={() => onSelect(p)}>
                {p.cover ? (
                  <img
                    src={p.cover}
                    alt={`${p.name} cover`}
                    className="w-[100px] h-[100px] object-cover rounded"
                    width={100}
                    height={100}
                  />
                ) : (
                  <div className="w-[100px] h-[100px] bg-gray-200 rounded" />
                )}
              </td>
              <td className="p-2 border cursor-pointer" onClick={() => onSelect(p)}>{p.name}</td>
              <td className="p-2 border cursor-pointer" onClick={() => onSelect(p)}>{truncate(p.description ?? '')}</td>
              <td className="p-2 border cursor-pointer" onClick={() => onSelect(p)}>{p.price}</td>
              <td className="p-2 border">
                <div className="flex items-center gap-2">
                  <Button
                    onClick={(e) => {
                      e.stopPropagation()
                      onUpdate(p)
                    }}
                  >
                    <AiFillEdit className="mr-1" /> Update
                  </Button>
                  <Button
                    onClick={(e) => {
                      e.stopPropagation()
                      handleDelete(p.id)
                    }}
                    className="bg-red-600 hover:bg-red-700"
                  >
                    <AiFillDelete className="mr-1" /> Delete
                  </Button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {loading && <div>Loading…</div>}
    </div>
  )
}
