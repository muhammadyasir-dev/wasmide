// Oauth.tsx
import React from 'react';
import { GoogleLogin } from 'react-google-login';

const Oauth: React.FC = () => {
  const responseGoogle = (response: any) => {
    console.log(response);
    // Send the token to your backend for verification
    fetch('http://localhost:8080/auth/google', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ token: response.tokenId }),
    })
      .then(res => res.json())
      .then(data => {
        console.log('Backend response:', data);
      })
      .catch(err => {
        console.error('Error:', err);
      });
  };

  const handleError = (error: any) => {
    console.error('Login failed:', error);
  };

  return (
    <div>
      <h2>Login with Google</h2>
      <GoogleLogin
        clientId="YOUR_GOOGLE_CLIENT_ID" // Replace with your Google client ID
        buttonText="Login with Google"
        onSuccess={responseGoogle}
        onFailure={handleError}
        cookiePolicy={'single_host_origin'}
      />
    </div>
  );
};

export default Oauth;