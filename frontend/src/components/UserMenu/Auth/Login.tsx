import { useState, useEffect } from 'react';
import * as React from 'react';
import { useAppDispatch, useAppSelector } from '../../../redux/hooks';
import { login, register, fetchGroupsByUserId } from '../../../redux/slices/authSlice';
import { RootState } from '../../../redux/store';
import './Login.css';


interface LoginProps {
  resetForm: boolean;
}


const Login: React.FC<LoginProps> = ({ resetForm }) => {
  const [username, setUsername] = useState(''); 
  const [password, setPassword] = useState('');
  const [firstname, setFirstname] = useState('');
  const [surname, setSurname] = useState('');
  const [role, setRole] = useState('student');
  const [isRegistering, setIsRegistering] = useState(false);
  const dispatch = useAppDispatch();
  const { user, loginLoading, loginError, registerLoading, registerError } = useAppSelector((state: RootState) => state.auth);

  useEffect(() => {
    if (resetForm) {
      setUsername('');
      setPassword('');
      setFirstname('');
      setSurname('');
      setRole('student');
    }
  }, [resetForm]);

  useEffect(() => {
    if (user) {
      dispatch(fetchGroupsByUserId());
    }
  }, [dispatch, user]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!username || !password || (isRegistering && (!firstname || !surname || !role))) {
      return;
    }

    try {
      if (isRegistering) {
        await dispatch(register({ username, password, firstname, surname, role })).unwrap();
      } else {
        await dispatch(login({ username, password })).unwrap();
        await dispatch(fetchGroupsByUserId()).unwrap();
      }
    } catch (err) {
      console.error(isRegistering ? 'Registration error:' : 'Login error:', err);
    }
  };

  return (
    <div className='auth-stuff'>
      {loginError && <p style={{ color: 'red' }}>{loginError}</p>}
      {registerError && <p style={{ color: 'red' }}>{registerError}</p>}
      {user ? (
        <div>
          <span className='username'>{user.surname} {user.firstname}</span>
        </div>
      ) : (
        <form onSubmit={handleSubmit}>
          <div className='auth-labels-outer'>
          <div className='auth-labels'>
          <div>
            <label>Логин: </label>
            <input type='text' value={username} onChange={(e) => setUsername(e.target.value)} />
          </div>
          <div>
            <label>Пароль: </label>
            <input type='password' value={password} onChange={(e) => setPassword(e.target.value)} />
          </div>
          {isRegistering && (
            <>
              <div>
                <label>Имя: </label>
                <input type='text' value={firstname} onChange={(e) => setFirstname(e.target.value)} />
              </div>
              <div>
                <label>Фамилия: </label>
                <input type='text' value={surname} onChange={(e) => setSurname(e.target.value)} />
              </div>
              <div>
                <label>Роль: </label>
                <select value={role} onChange={(e) => setRole(e.target.value)}>
                  <option value='student'>Студент</option>
                  <option value='teacher'>Учитель</option>
                </select>
              </div>
            </>
          )}
          </div>
          </div>
          <div className="button-container-auth">
            <button type='submit' disabled={loginLoading || registerLoading}>
              {loginLoading || registerLoading ? 'Загрузка...' : isRegistering ? 'Зарегистрироваться' : 'Войти'}
            </button>
            <button type='button' onClick={() => setIsRegistering(!isRegistering)}>
              {isRegistering ? 'Уже есть аккаунт? Войти' : 'Нет аккаунта? Зарегистрироваться'}
            </button>
          </div>
        </form>
      )}
    </div>
  );
};


export default Login;
