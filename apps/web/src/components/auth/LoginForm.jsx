import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { login } from '../../services/auth';

export default function LoginForm() {
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [fieldErrors, setFieldErrors] = useState({});

  function validateEmail(value) {
    if (!value) return 'This field is required.';
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(value)) return 'Enter a valid email address.';
    return '';
  }

  function validatePassword(value) {
    if (!value) return 'Password is required';
    return '';
  }

  function handleSubmit(e) {
    e.preventDefault();
    setError('');
    setFieldErrors({});

    const emailError = validateEmail(email);
    const passwordError = validatePassword(password);

    if (emailError || passwordError) {
      setFieldErrors({ email: emailError, password: passwordError });
      return;
    }

    setLoading(true);

    login(email, password)
      .then(() => {
        navigate('/dashboard');
      })
      .catch((err) => {
        setError(err.message);
      })
      .finally(() => {
        setLoading(false);
      });
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6" noValidate>
      {error && (
        <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
          {error}
        </div>
      )}

      <div className="space-y-2">
        <label htmlFor="email" className="block text-base font-medium text-gray-700">
          Email
        </label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
            fieldErrors.email
              ? 'border-red-600 focus:border-red-600'
              : 'border-gray-200 focus:border-blue-600'
          }`}
          placeholder="you@example.com"
          autoComplete="email"
          aria-describedby={fieldErrors.email ? 'email-error' : undefined}
          aria-invalid={!!fieldErrors.email}
        />
        {fieldErrors.email && (
          <p id="email-error" className="text-sm text-red-600" role="alert">
            {fieldErrors.email}
          </p>
        )}
      </div>

      <div className="space-y-2">
        <label htmlFor="password" className="block text-base font-medium text-gray-700">
          Password
        </label>
        <input
          id="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
            fieldErrors.password
              ? 'border-red-600 focus:border-red-600'
              : 'border-gray-200 focus:border-blue-600'
          }`}
          placeholder="Enter your password"
          autoComplete="current-password"
          aria-describedby={fieldErrors.password ? 'password-error' : undefined}
          aria-invalid={!!fieldErrors.password}
        />
        {fieldErrors.password && (
          <p id="password-error" className="text-sm text-red-600" role="alert">
            {fieldErrors.password}
          </p>
        )}
      </div>

      <div className="flex justify-end">
        <Link to="/reset" className="text-sm text-blue-600 hover:text-blue-700">
          Forgot password?
        </Link>
      </div>

      <button
        type="submit"
        disabled={loading}
        className="min-h-[44px] w-full rounded-md bg-blue-600 px-4 py-2 text-base font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
      >
        {loading ? 'Signing in...' : 'Sign In'}
      </button>

      <p className="text-center text-sm text-gray-600">
        Don&apos;t have an account?{' '}
        <Link to="/register" className="text-blue-600 hover:text-blue-700">
          Create one
        </Link>
      </p>
    </form>
  );
}
