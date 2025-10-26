import { Routes, Route, Navigate, Link } from 'react-router-dom'
import ListPage from './pages/ListPage'
import DetailPage from './pages/DetailPage'
import CreatePage from './pages/CreatePage'
import UpdatePage from './pages/UpdatePage'
import {UserCreatePage} from "./pages/user"

export default function App() {
  return (
    <div className="min-h-screen">
      <header className="border-b p-4">
        <div className="mx-auto max-w-5xl flex items-center justify-between">
          <Link to="/" className="text-lg font-semibold">Atlas</Link>
        </div>
      </header>
      <main className="mx-auto max-w-5xl">
        <Routes>
          <Route path="/" element={<ListPage />} />
          <Route path="/products/new" element={<CreatePage />} />
          <Route path="/products/:id" element={<DetailPage />} />
          <Route path="/products/:id/edit" element={<UpdatePage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
          <Route path="/users/new" element={<UserCreatePage />} />
        </Routes>
      </main>
    </div>
  )
}
