import * as React from 'react';
import './UserMenu.css';
import { useNavigate } from 'react-router-dom';
import { LogoutIcon } from './LogoutIcon';
import { AvatarIcon } from './AvatarIcon';
import Login from './Auth/Login';
import { useAppDispatch } from '../../redux/hooks';
import { logout } from '../../redux/slices/authSlice';
import { clearContests } from '../../redux/slices/contestsSlice';


const UserMenu: React.FC = () => {
  const dispatch = useAppDispatch();
  const [resetLoginForm, setResetLoginForm] = React.useState(false);
  const navigate = useNavigate();

  const handleLogout = () => {
    dispatch(logout());
    dispatch(clearContests());
    setResetLoginForm(true);

    setTimeout(() => {
      setResetLoginForm(false);
    }, 100);

    navigate('');
  };

  return (
    <div className='user-menu'>
      <Login resetForm={resetLoginForm} />
      <div className='avatar'>  
        <AvatarIcon />
      </div>
      <button className='logout-btn' onClick={handleLogout}>
        <LogoutIcon />
      </button>
    </div>
  );
};


export default UserMenu;
