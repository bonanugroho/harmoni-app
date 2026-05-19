import { useState, type FormEvent, type ChangeEvent } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { register } from '../../services/auth';

interface Territory {
  id: string;
  name: string;
}

interface PasswordStrength {
  level: string;
  label: string;
  color: string;
}

const TERRITORIES: Territory[] = [
  { id: 'rt-01', name: 'RT 01' },
  { id: 'rt-02', name: 'RT 02' },
  { id: 'rw-01', name: 'RW 01' },
];

function validatePassword(password: string): string[] {
  const errors: string[] = [];
  if (password.length < 8) {
    errors.push('Password must be at least 8 characters.');
  }
  if (!/[A-Z]/.test(password)) {
    errors.push('Password must include uppercase, lowercase, numbers, and symbols.');
  }
  if (!/[a-z]/.test(password)) {
    errors.push('Password must include uppercase, lowercase, numbers, and symbols.');
  }
  if (!/[0-9]/.test(password)) {
    errors.push('Password must include uppercase, lowercase, numbers, and symbols.');
  }
  if (!/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
    errors.push('Password must include uppercase, lowercase, numbers, and symbols.');
  }
  return errors;
}

function getPasswordStrength(password: string): PasswordStrength {
  if (!password) return { level: '', label: '', color: '' };

  let score = 0;
  if (password.length >= 8) score++;
  if (/[A-Z]/.test(password)) score++;
  if (/[a-z]/.test(password)) score++;
  if (/[0-9]/.test(password)) score++;
  if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) score++;

  if (score <= 2) return { level: 'weak', label: 'Weak', color: 'text-red-600' };
  if (score <= 3) return { level: 'medium', label: 'Medium', color: 'text-yellow-600' };
  return { level: 'strong', label: 'Strong', color: 'text-green-600' };
}

export default function RegisterForm() {
  const navigate = useNavigate();
  const [fullName, setFullName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [territoryId, setTerritoryId] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  function validateEmail(value: string): string {
    if (!value) return 'This field is required.';
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(value)) return 'Enter a valid email address.';
    return '';
  }

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setError('');
    setFieldErrors({});

    const errors: Record<string, string> = {};
    if (!fullName.trim()) errors.fullName = 'This field is required.';
    const emailError = validateEmail(email);
    if (emailError) errors.email = emailError;
    const passwordErrors = validatePassword(password);
    if (passwordErrors.length > 0) errors.password = passwordErrors[0];
    if (password !== confirmPassword) errors.confirmPassword = 'Passwords do not match.';
    if (!territoryId) errors.territoryId = 'This field is required.';

    if (Object.keys(errors).length > 0) {
      setFieldErrors(errors);
      return;
    }

    setLoading(true);

    register({ email, password, fullName, territoryId })
      .then(() => {
        navigate('/login');
      })
      .catch((err: Error) => {
        if (err.message.includes('already registered')) {
          setFieldErrors({ email: 'An account with this email already exists. Try signing in instead.' });
        } else {
          setError(err.message);
        }
      })
      .finally(() => {
        setLoading(false);
      });
  }

  const strength = getPasswordStrength(password);

  return (
    <form onSubmit={handleSubmit} className="space-y-5" noValidate>
      {error && (
        <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
          {error}
        </div>
      )}

      <div className="space-y-2">
        <label htmlFor="fullName" className="block text-base font-medium text-gray-700">
          Full Name
        </label>
        <input
          id="fullName"
          type="text"
          value={fullName}
          onChange={(e: ChangeEvent<HTMLInputElement>) => setFullName(e.target.value)}
          className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
            fieldErrors.fullName
              ? 'border-red-600 focus:border-red-600'
              : 'border-gray-200 focus:border-blue-600'
          }`}
          placeholder="John Doe"
          autoComplete="name"
          aria-invalid={!!fieldErrors.fullName}
        />
        {fieldErrors.fullName && (
          <p className="text-sm text-red-600" role="alert">
            {fieldErrors.fullName}
          </p>
        )}
      </div>

      <div className="space-y-2">
        <label htmlFor="email" className="block text-base font-medium text-gray-700">
          Email
        </label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={(e: ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
          className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
            fieldErrors.email
              ? 'border-red-600 focus:border-red-600'
              : 'border-gray-200 focus:border-blue-600'
          }`}
          placeholder="you@example.com"
          autoComplete="email"
          aria-invalid={!!fieldErrors.email}
        />
        {fieldErrors.email && (
          <p className="text-sm text-red-600" role="alert">
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
          onChange={(e: ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
          className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
            fieldErrors.password
              ? 'border-red-600 focus:border-red-600'
              : 'border-gray-200 focus:border-blue-600'
          }`}
          placeholder="Create a strong password"
          autoComplete="new-password"
          aria-invalid={!!fieldErrors.password}
        />
        {password && (
          <p className={`text-sm font-medium ${strength.color}`}>
            Password strength: {strength.label}
          </p>
        )}
        <p className="text-xs text-gray-500">
          Use 8+ characters with uppercase, lowercase, numbers, and symbols.
        </p>
        {fieldErrors.password && (
          <p className="text-sm text-red-600" role="alert">
            {fieldErrors.password}
          </p>
        )}
      </div>

      <div className="space-y-2">
        <label htmlFor="confirmPassword" className="block text-base font-medium text-gray-700">
          Confirm Password
        </label>
        <input
          id="confirmPassword"
          type="password"
          value={confirmPassword}
          onChange={(e: ChangeEvent<HTMLInputElement>) => setConfirmPassword(e.target.value)}
          className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
            fieldErrors.confirmPassword
              ? 'border-red-600 focus:border-red-600'
              : 'border-gray-200 focus:border-blue-600'
          }`}
          placeholder="Confirm your password"
          autoComplete="new-password"
          aria-invalid={!!fieldErrors.confirmPassword}
        />
        {fieldErrors.confirmPassword && (
          <p className="text-sm text-red-600" role="alert">
            {fieldErrors.confirmPassword}
          </p>
        )}
      </div>

      <div className="space-y-2">
        <label htmlFor="territory" className="block text-base font-medium text-gray-700">
          Territory
        </label>
        <select
          id="territory"
          value={territoryId}
          onChange={(e: ChangeEvent<HTMLSelectElement>) => setTerritoryId(e.target.value)}
          className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
            fieldErrors.territoryId
              ? 'border-red-600 focus:border-red-600'
              : 'border-gray-200 focus:border-blue-600'
          }`}
          aria-invalid={!!fieldErrors.territoryId}
        >
          <option value="">Select your territory</option>
          {TERRITORIES.map((t) => (
            <option key={t.id} value={t.id}>
              {t.name}
            </option>
          ))}
        </select>
        {fieldErrors.territoryId && (
          <p className="text-sm text-red-600" role="alert">
            {fieldErrors.territoryId}
          </p>
        )}
      </div>

      <button
        type="submit"
        disabled={loading}
        className="min-h-[44px] w-full rounded-md bg-blue-600 px-4 py-2 text-base font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
      >
        {loading ? 'Creating account...' : 'Create Account'}
      </button>

      <p className="text-center text-sm text-gray-600">
        Already have an account?{' '}
        <Link to="/login" className="text-blue-600 hover:text-blue-700">
          Sign in
        </Link>
      </p>
    </form>
  );
}
