import { Routes, Route } from 'react-router-dom';
import { Layout } from '@/components/Layout';
import { HomePage } from '@/pages/HomePage';
import GeneratePage from '@/pages/GeneratePage';
import ForumPage from '@/pages/ForumPage';
import ProfilePage from '@/pages/ProfilePage';
import AdminPage from '@/pages/AdminPage';
import GalleryPage from '@/pages/GalleryPage';
import HistoryPage from '@/pages/HistoryPage';
import PromptsPage from '@/pages/PromptsPage';

export function AppRouter() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="generate" element={<GeneratePage />} />
        <Route path="forum" element={<ForumPage />} />
        <Route path="profile" element={<ProfilePage />} />
        <Route path="admin" element={<AdminPage />} />
        <Route path="gallery" element={<GalleryPage />} />
        <Route path="history" element={<HistoryPage />} />
        <Route path="prompts" element={<PromptsPage />} />
      </Route>
    </Routes>
  );
}