import { useState } from 'react'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { Button } from './ui/button'

export interface ProductFormValues {
  id?: number
  name: string
  description: string
  price: number
  cover: string
}

interface Props {
  initial?: ProductFormValues
  onSubmit: (values: ProductFormValues) => void
}

export function ProductForm({ initial, onSubmit }: Props) {
  const [name, setName] = useState(initial?.name ?? '')
  const [description, setDescription] = useState(initial?.description ?? '')
  const [price, setPrice] = useState(initial?.price ?? 0)
  const [cover, setCover] = useState(initial?.cover ?? '')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit({ id: initial?.id, name, description, price, cover })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Name</Label>
        <Input id="name" value={name} onChange={e => setName(e.target.value)} />
      </div>
      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <textarea
          id="description"
          className="block w-full rounded border border-gray-300 bg-white px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-blue-500"
          rows={4}
          value={description}
          onChange={e => setDescription(e.target.value)}
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="price">Price</Label>
        <Input id="price" type="number" value={price} onChange={e => setPrice(parseFloat(e.target.value))} />
      </div>
      <div className="space-y-2">
        <Label htmlFor="cover">Cover URL</Label>
        <Input id="cover" value={cover} onChange={e => setCover(e.target.value)} />
      </div>
      <Button type="submit">Save</Button>
    </form>
  )
}
