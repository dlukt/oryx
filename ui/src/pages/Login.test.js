import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import Login from './Login';
import { BrowserRouter } from 'react-router-dom';
import { I18nextProvider } from 'react-i18next';
import axios from 'axios';
import '@testing-library/jest-dom';

// Mock react-i18next
jest.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key) => key }),
  I18nextProvider: ({ children }) => <div>{children}</div>,
  initReactI18next: { type: '3rdParty', init: jest.fn() },
}));

// Mock react-error-boundary
jest.mock('react-error-boundary', () => ({
  useErrorBoundary: () => ({ showBoundary: jest.fn() }),
  ErrorBoundary: ({ children }) => <div>{children}</div>,
}));

// Mock axios
jest.mock('axios');

// Mock utils and i18n
jest.mock('../utils', () => ({
  Token: {
    load: jest.fn(),
    loadBearerHeader: jest.fn(),
    save: jest.fn(),
  },
  Tools: {
    mask: jest.fn(),
  },
  Locale: {
      load: jest.fn(),
      current: jest.fn(() => 'en'),
  }
}));

// Mock i18n module
jest.mock('../i18n', () => ({
    use: () => ({
        init: jest.fn()
    })
}));

// Mock SrsErrorBoundary
jest.mock('../components/SrsErrorBoundary', () => ({
  SrsErrorBoundary: ({ children }) => <div>{children}</div>,
}));

// Mock useNavigate
const mockedNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockedNavigate,
}));

describe('Login Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders login form', () => {
    render(
      <BrowserRouter>
        <Login />
      </BrowserRouter>
    );

    expect(screen.getByText('login.passwordLabel')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Password')).toBeInTheDocument();
    expect(screen.getByText('login.labelLogin')).toBeInTheDocument();
  });

  test('toggle password visibility', () => {
    render(
      <BrowserRouter>
        <Login />
      </BrowserRouter>
    );

    // Initial state: plaintext is true, so "Hide password" should be the label
    const toggleButton = screen.getByLabelText('Hide password');
    const passwordInput = screen.getByPlaceholderText('Password');

    // Default state: plaintext is true in the component (based on initial state true)
    // So initially it should be type="text".
    expect(passwordInput).toHaveAttribute('type', 'text');

    // Click toggle button
    fireEvent.click(toggleButton);

    // Now it should be password type
    expect(passwordInput).toHaveAttribute('type', 'password');

    // And button label should change
    expect(screen.getByLabelText('Show password')).toBeInTheDocument();
  });
});
