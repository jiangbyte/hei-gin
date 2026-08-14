/** Author: Charlie */

import { Navigate, type RouteObject } from 'react-router-dom'
import { MainLayout } from '@/layouts'
import { HomePage } from '@/pages/home'
import { LoginPage } from '@/pages/auth/login'
import { RegisterPage } from '@/pages/auth/register'
import { ForgotPasswordPage } from '@/pages/auth/forgot-password'
import { OAuthCallbackPage } from '@/pages/auth/oauth-callback'
import { NotFoundPage } from '@/pages/error/not-found'
import { UserCenterPage } from '@/pages/usercenter'
import { ProfilePage } from '@/pages/profile'
import { AnnouncementListPage } from '@/pages/announcements'
import { AnnouncementDetailPage } from '@/pages/announcements/detail'
import { FeedbackListPage } from '@/pages/feedback'
import { FeedbackNewPage } from '@/pages/feedback/new'
import { FeedbackDetailPage } from '@/pages/feedback/detail'
import { guestOnly, requireAuth } from './guard'

export const routes: RouteObject[] = [
  {
    path: '/',
    element: <MainLayout />,
    children: [
      { index: true, element: <HomePage /> },
      {
        path: 'profile',
        element: <ProfilePage />,
      },
      {
        path: 'usercenter',
        loader: requireAuth,
        element: <UserCenterPage />,
      },
      {
        path: 'announcements',
        element: <AnnouncementListPage />,
      },
      {
        path: 'announcements/:id',
        element: <AnnouncementDetailPage />,
      },
      {
        path: 'feedback',
        loader: requireAuth,
        element: <FeedbackListPage />,
      },
      {
        path: 'feedback/new',
        loader: requireAuth,
        element: <FeedbackNewPage />,
      },
      {
        path: 'feedback/:id',
        loader: requireAuth,
        element: <FeedbackDetailPage />,
      },
    ],
  },
  {
    path: '/auth/login',
    loader: guestOnly,
    element: <LoginPage />,
  },
  {
    path: '/auth/register',
    loader: guestOnly,
    element: <RegisterPage />,
  },
  {
    path: '/auth/forgot-password',
    loader: guestOnly,
    element: <ForgotPasswordPage />,
  },
  {
    path: '/auth/oauth/callback',
    element: <OAuthCallbackPage />,
  },
  { path: '/404', element: <NotFoundPage /> },
  { path: '*', element: <Navigate to="/404" replace /> },
]
