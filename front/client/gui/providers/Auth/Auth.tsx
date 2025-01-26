import React, {useEffect, useState} from 'react';
import api from '../../api';
import i18n from '../../i18n';

const Auth = props => {
  const [authed, setAuthed] = useState(false)
  const [loading, setLoading] = useState(true)
  const [_, setError] = useState(null)

  useEffect(() => {
    async function call() {
      setLoading(true)
      try {
        let response = await api.auth()
        if (response.error == "") {
          setAuthed(true)
        } else {
          setError(response.error)
        }
      } catch (e) {
        setError(e)
      } finally {
        setLoading(false)
      }
    }

    call();
  }, [setAuthed, setLoading, setError])

  if (loading) {
    return (<>{i18n('auth_process')}</>)
  }

  return (
    <>{authed
      ? props.children
      : (props.fallback || i18n("not_authed"))
    }</>
  );
}

export default Auth;
