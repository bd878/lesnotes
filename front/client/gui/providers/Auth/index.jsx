import React, {useEffect, useState} from 'react';
import api from '../../api';

const AuthProvider = props => {
  const [authed, setAuthed] = useState(false)
  useEffect(async () => {
    const response = await api("/users/v1/auth", {
      method: 'POST',
      credentials: 'include',
    });
    console.log(response)
  }, [])

  return (
    <>{authed
      ? props.children
      : (props.fallback || "Authenticating")
    }</>
  );
}

export default AuthProvider;
