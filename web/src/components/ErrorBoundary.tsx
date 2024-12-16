import React, { Component, ErrorInfo, ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';

interface ErrorBoundaryProps {
  children: ReactNode;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
}

class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { 
      hasError: false,
      error: null,
      errorInfo: null 
    };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error: error, errorInfo: null };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    const navigate = useNavigate(); // Usamos el navigate aquí.
    console.error('ErrorBoundary caught an error', error, errorInfo);

    // Redirige al login con el mensaje de error
    navigate('/login', { state: { errorMessage: error.message || 'Algo salió mal.' } });
  }

  render() {
    if (this.state.hasError) {
      return null; // No renderizamos nada aquí, ya que redirigimos al login
    }

    return this.props.children; 
  }
}

export default ErrorBoundary;
