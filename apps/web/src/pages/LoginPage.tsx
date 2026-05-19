import LoginForm from '../components/auth/LoginForm';

export default function LoginPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-surface px-4 py-8 sm:px-6">
      <div className="w-full max-w-[400px]">
        <div className="mb-8 text-center sm:mb-12">
          <h1 className="text-[28px] font-semibold text-text-primary">Harmoni</h1>
          <p className="mt-2 text-sm text-text-secondary">
            Community Financial Management
          </p>
        </div>

        <div className="rounded-lg bg-white p-6 shadow-sm sm:p-8">
          <h2 className="mb-6 text-center text-[24px] font-semibold text-text-primary">
            Sign In to Harmoni
          </h2>
          <LoginForm />
        </div>

        <footer className="mt-8 text-center text-xs text-text-muted">
          &copy; {new Date().getFullYear()} Harmoni. All rights reserved.
        </footer>
      </div>
    </div>
  );
}
