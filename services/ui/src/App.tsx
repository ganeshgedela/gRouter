import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { QueryClientProvider } from "@tanstack/react-query"
import { ThemeProvider } from "@/components/theme-provider"
import { MainLayout } from "@/components/layout/main-layout"
import { queryClient } from "@/lib/query-client"

// Placeholder pages (we will implement these next)
import Dashboard from "@/pages/dashboard"

import Tenants from "@/pages/tenants"
import Users from "@/pages/users" // Re-using old name/file if I revert renaming too? The file was renamed to iam.tsx. I should check if I should revert the file rename.
import Network from "@/pages/network"
import Metrics from "@/pages/metrics"
import Alerts from "@/pages/alerts"
import Chat from "@/pages/chat"
import Settings from "@/pages/settings"
import TenantDetail from "@/pages/tenant-detail"
import RouterDetail from "@/pages/router-detail"
import PoliciesPage from "@/pages/policies"
import SecurityPage from "@/pages/security"
import { Toaster } from "@/components/ui/toaster"
import { GlobalLoading } from "@/components/ui/global-loading"
import { ErrorBoundary } from "@/components/error-boundary"

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
        <BrowserRouter>
          <ErrorBoundary>
            <GlobalLoading />
            <Routes>
              <Route element={<MainLayout />}>
                <Route path="/" element={<Dashboard />} />
                <Route path="/tenants" element={<Tenants />} />
                <Route path="/tenants/:id" element={<TenantDetail />} />
                <Route path="/users" element={<Users />} />
                <Route path="/routers" element={<Network />} />
                <Route path="/routers/:id" element={<RouterDetail />} />
                <Route path="/policies" element={<PoliciesPage />} />
                <Route path="/security" element={<SecurityPage />} />
                <Route path="/metrics" element={<Metrics />} />
                <Route path="/alerts" element={<Alerts />} />
                <Route path="/chat" element={<Chat />} />
                <Route path="/settings" element={<Settings />} />
                <Route path="*" element={<Navigate to="/" replace />} />
              </Route>
            </Routes>
            <Toaster />
          </ErrorBoundary>
        </BrowserRouter>
      </ThemeProvider>
    </QueryClientProvider>
  )
}

export default App
