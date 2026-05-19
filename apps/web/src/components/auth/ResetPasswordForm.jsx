import { useState, useEffect } from 'react';
import { useSearchParams, useNavigate, Link } from 'react-router-dom';
import { requestPasswordReset, confirmPasswordReset } from '../../services/auth';

function validatePassword(password) {
  const errors = [];
  if (password.length < 8) {
    errors.push('Password must be at least 8 characters.');
  }
  if (!/[A-Z]/.test(password) || !/[a-z]/.test(password) || !/[0-9]/.test(password) || !/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
    errors.push('Password must include uppercase, lowercase, numbers, and symbols.');
  }
  return errors;
}

function getPasswordStrength(password) {
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

export default function ResetPasswordForm() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get('token');

  const [email, setEmail] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [fieldErrors, setFieldErrors] = useState({});

  // If token is present in URL, show new password form
  const isNewPasswordMode = !!token;

  useEffect(() => {
    if (token) {
      setNewPassword('');
      setConfirmPassword('');
    }
  }, [token]);

  function handleRequestReset(e) {
    e.preventDefault();
    setError('');
    setSuccess('');
    setFieldErrors({});

    if (!email) {
      setFieldErrors({ email: 'This field is required.' });
      return;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setFieldErrors({ email: 'Enter a valid email address.' });
      return;
    }

    setLoading(true);

    requestPasswordReset(email)
      .then(() => {
        setSuccess(`Check your email. We sent a reset link to ${email}. It expires in 1 hour.`);
        setEmail('');
      })
      .catch((err) => {
        setError(err.message);
      })
      .finally(() => {
        setLoading(false);
      });
  }

  function handleConfirmReset(e) {
    e.preventDefault();
    setError('');
    setSuccess('');
    setFieldErrors({});

    const errors = {};
    const passwordErrors = validatePassword(newPassword);
    if (passwordErrors.length > 0) errors.newPassword = passwordErrors[0];
    if (newPassword !== confirmPassword) errors.confirmPassword = 'Passwords do not match.';

    if (Object.keys(errors).length > 0) {
      setFieldErrors(errors);
      return;
    }

    setLoading(true);

    confirmPasswordReset(token, newPassword)
      .then(() => {
        setSuccess('Password updated. You can now sign in with your new password.');
        setTimeout(() => navigate('/login'), 2000);
      })
      .catch((err) => {
        if (err.message.includes('Invalid') || err.message.includes('expired')) {
          setError('Invalid or expired reset token');
        } else {
          setError(err.message);
        }
      })
      .finally(() => {
        setLoading(false);
      });
  }

  const strength = getPasswordStrength(newPassword);

  // New password mode (token present)
  if (isNewPasswordMode) {
    return (
      <form onSubmit={handleConfirmReset} className="space-y-5" noValidate>
        {error && (
          <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
            {error}
          </div>
        )}
        {success && (
          <div role="status" className="rounded-md bg-green-50 p-4 text-sm text-green-700">
            {success}
          </div>
        )}

        <div className="space-y-2">
          <label htmlFor="newPassword" className="block text-base font-medium text-gray-700">
            New Password
          </label>
          <input
            id="newPassword"
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
              fieldErrors.newPassword
                ? 'border-red-600 focus:border-red-600'
                : 'border-gray-200 focus:border-blue-600'
            }`}
            placeholder="Create a strong password"
            autoComplete="new-password"
            aria-invalid={!!fieldErrors.newPassword}
          />
          {newPassword && (
            <p className={`text-sm font-medium ${strength.color}`}>
              Password strength: {strength.label}
            </p>
          )}
          <p className="text-xs text-gray-500">
            Use 8+ characters with uppercase, lowercase, numbers, and symbols.
          </p>
          {fieldErrors.newPassword && (
            <p className="text-sm text-red-600" role="alert">
              {fieldErrors.newPassword}
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
            onChange={(e) => setConfirmPassword(e.target.value)}
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

        <button
          type="submit"
          disabled={loading}
          className="min-h-[44px] w-full rounded-md bg-blue-600 px-4 py-2 text-base font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {loading ? 'Updating password...' : 'Update Password'}
        </button>

        <p className="text-center text-sm text-gray-600">
          <Link to="/login" className="text-blue-600 hover:text-blue-700">
            Back to Sign In
          </Link>
        </p>
      </form>
    );
  }

  // Request reset mode (no token)
  return (
    <form onSubmit={handleRequestReset} className="space-y-5" noValidate>
      {error && (
        <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
          {error}
        </div>
      )}
      {success && (
        <div role="status" className="rounded-md bg-green-50 p-4 text-sm text-green-700">
          {success}
        </div>
      )}

      {!success && (
        <>
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
              aria-invalid={!!fieldErrors.email}
            />
            {fieldErrors.email && (
              <p className="text-sm text-red-600" role="alert">
                {fieldErrors.email}
              </p>
            )}
          </div>

          <button
            type="submit"
            disabled={loading}
            className="min-h-[44px] w-full rounded-md bg-blue-600 px-4 py-2 text-base font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {loading ? 'Sending...' : 'Send Reset Link'}
          </button>
        </>
      )}

      <p className="text-center text-sm text-gray-600">
        <Link to="/login" className="text-blue-600 hover:text-blue-700">
          Back to Sign In
        </Link>
      </p>
    </form>
  );
}
