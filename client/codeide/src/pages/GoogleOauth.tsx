import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { backendurl } from "./../libs/Url";
// Define our interface for the user
interface User {
  id: string;
  name: string;
  email: string;
  picture?: string;
}

const Oauth: React.FC = () => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const navigate = useNavigate();

  // Check if user is already logged in on component mount
  useEffect(() => {
    const checkAuthStatus = async () => {
      try {
        const response = await axios.get("${backendurl}api/auth/user", {
          withCredentials: true,
        });
        if (response.data.user) {
          setUser(response.data.user);
        }
      } catch (error) {
        console.error("Failed to fetch user:", error);
      } finally {
        setLoading(false);
      }
    };

    checkAuthStatus();
  }, []);

  // Handle Google login
  const handleGoogleLogin = () => {
    window.location.href = `${backendurl}api/auth/google/login`;
  };

  // Handle logout
  const handleLogout = async () => {
    try {
      await axios.post(`${backendurl}api/auth/logout`, {}, {
        withCredentials: true,
      });
      setUser(null);
      navigate("/");
    } catch (error) {
      console.error("Failed to logout:", error);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-screen">Loading</div>
    );
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 p-4">
      {user
        ? (
          <div className="bg-white p-8 rounded-lg shadow-md max-w-md w-full">
            <div className="flex flex-col items-center">
              {user.picture && (
                <img
                  src={user.picture}
                  alt={user.name}
                  className="w-20 h-20 rounded-full mb-4"
                />
              )}
              <h2 className="text-2xl font-bold mb-2">{user.name}</h2>
              <p className="text-gray-600 mb-6">{user.email}</p>
              <button
                onClick={handleLogout}
                className="w-full bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded-md transition duration-200"
              >
                Logout
              </button>
            </div>
          </div>
        )
        : (
          <div className="bg-white p-8 rounded-lg shadow-md max-w-md w-full">
            <h2 className="text-2xl font-bold mb-6 text-center">Sign in</h2>
            <button
              onClick={handleGoogleLogin}
              className="w-full flex items-center justify-center bg-white border border-gray-300 rounded-md py-2 px-4 hover:bg-gray-50 transition duration-200"
            >
              <svg
                viewBox="0 0 24 24"
                className="w-6 h-6 mr-2"
                xmlns="http://www.w3.org/2000/svg"
              >
                <g transform="matrix(1, 0, 0, 1, 27.009001, -39.238998)">
                  <path
                    fill="#4285F4"
                    d="M -3.264 51.509 C -3.264 50.719 -3.334 49.969 -3.454 49.239 L -14.754 49.239 L -14.754 53.749 L -8.284 53.749 C -8.574 55.229 -9.424 56.479 -10.684 57.329 L -10.684 60.329 L -6.824 60.329 C -4.564 58.239 -3.264 55.159 -3.264 51.509 Z"
                  />
                  <path
                    fill="#34A853"
                    d="M -14.754 63.239 C -11.514 63.239 -8.804 62.159 -6.824 60.329 L -10.684 57.329 C -11.764 58.049 -13.134 58.489 -14.754 58.489 C -17.884 58.489 -20.534 56.379 -21.484 53.529 L -25.464 53.529 L -25.464 56.619 C -23.494 60.539 -19.444 63.239 -14.754 63.239 Z"
                  />
                  <path
                    fill="#FBBC05"
                    d="M -21.484 53.529 C -21.734 52.809 -21.864 52.039 -21.864 51.239 C -21.864 50.439 -21.724 49.669 -21.484 48.949 L -21.484 45.859 L -25.464 45.859 C -26.284 47.479 -26.754 49.299 -26.754 51.239 C -26.754 53.179 -26.284 54.999 -25.464 56.619 L -21.484 53.529 Z"
                  />
                  <path
                    fill="#EA4335"
                    d="M -14.754 43.989 C -12.984 43.989 -11.404 44.599 -10.154 45.789 L -6.734 42.369 C -8.804 40.429 -11.514 39.239 -14.754 39.239 C -19.444 39.239 -23.494 41.939 -25.464 45.859 L -21.484 48.949 C -20.534 46.099 -17.884 43.989 -14.754 43.989 Z"
                  />
                </g>
              </svg>
              Sign in with Google
            </button>
          </div>
        )}
    </div>
  );
};

export default Oauth;
