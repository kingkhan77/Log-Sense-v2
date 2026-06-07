import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Layout from './Layout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Alerts from './pages/Alerts';
import Rules from './pages/Rules';
import Services from './pages/Services';
import Users from './pages/Users';
import Logs from './pages/Logs';
import OnCall from './pages/OnCall';
import APIKeys from './pages/APIKeys';
import Profile from './pages/Profile';

function Guard({ children }) {
  return localStorage.getItem('token') ? children : <Navigate to="/login" replace />;
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={<Guard><Layout /></Guard>}>
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="alerts"    element={<Alerts />} />
          <Route path="rules"     element={<Rules />} />
          <Route path="services"  element={<Services />} />
          <Route path="users"     element={<Users />} />
          <Route path="logs"      element={<Logs />} />
          <Route path="oncall"    element={<OnCall />} />
          <Route path="api-keys"  element={<APIKeys />} />
          <Route path="profile"   element={<Profile />} />
        </Route>
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
