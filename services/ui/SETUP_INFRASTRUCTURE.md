# Production Infrastructure Setup

## Infrastructure Overview

This project is set up with:
- **React Query** for server state management
- **Vite** for build tooling
- **TypeScript** for type safety
- **Tailwind CSS** for styling

## Files Created

### Environment Configuration
- ✅ `.env.development` - Development environment variables
- ✅ `.env.production` - Production environment variables  
- ✅ `.env.example` - Template file for reference
- ✅ `src/lib/env.ts` - Typed environment helper

### API & State Management
- ✅ `src/lib/api-client.ts` - Enhanced with exponential backoff retry
- ✅ `src/lib/query-client.ts` - React Query configuration
- ✅ `src/hooks/use-api.ts` - Example custom hooks for data fetching

### App Integration
- ✅ `src/App.tsx` - Integrated QueryClientProvider

## Next Steps

1. Install the dependencies listed above
2. Update `package.json` to include the new dependencies  
3. Test the build with:
   ```bash
   npm run build
   ```
4. The mock data in `use-api.ts` can be replaced with real API calls once backend is ready

## Usage Example

### In your components:

```tsx
import { useRouters } from '@/hooks/use-api'

function RoutersList() {
  const { data: routers, isLoading, error } = useRouters()
  
  if (isLoading) return <div>Loading...</div>
  if (error) return <div>Error: {error.message}</div>
  
  return (
    <ul>
      {routers?.map(router => (
        <li key={router.id}>{router.hostname}</li>
      ))}
    </ul>
  )
}
```

## Environment Variables

Configure these in your `.env.development` and `.env.production`:

- `VITE_API_URL` - Backend API URL
- `VITE_WS_URL` - WebSocket URL for real-time updates
- `VITE_ENABLE_MOCKS` - Enable/disable mock data
- `VITE_SENTRY_DSN` - Sentry error tracking DSN (optional)
