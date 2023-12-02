import React, {useEffect, useState} from 'react';
import api from '../../api';

const AuthProvider = props => {
  const [authed, setAuthed] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function call() {
      let response = { valid: false };
      try {
        setLoading(true);
        response = await api("/users/v1/auth", {
          method: 'POST',
          credentials: 'include',
        });
      } catch(e) {
        console.error("error occured on authing:", e);
      } finally {
        setLoading(false);
      }

      if (response.valid) {
        setAuthed(true);
        console.log("welcome,", response.user.name)
      } else {
        setAuthed(false);
      }
    }

    call();
  }, [setAuthed, setLoading])

  if (loading) {
    return (<>Authenticating...</>)
  }

  return (
    <>{authed
      ? props.children
      : (props.fallback || "Not authed")
    }</>
  );
}

export default AuthProvider;
