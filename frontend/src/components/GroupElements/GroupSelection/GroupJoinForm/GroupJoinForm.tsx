import * as React from 'react';
import { useState } from 'react';
import { useAppDispatch, useAppSelector } from '../../../../redux/hooks';
import { joinGroupByCode } from '../../../../redux/slices/authSlice';
import './GroupJoinForm.css';


interface GroupJoinFormProps {
  closeForm: () => void;
}


const GroupJoinForm: React.FC<GroupJoinFormProps> = ({ closeForm }) => {
  const dispatch = useAppDispatch();
  const [groupCode, setGroupCode] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const { joinGroupByCodeLoading, joinGroupByCodeError } = useAppSelector((state) => state.auth);

  const handleFormSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!groupCode) {
      setError('Код группы обязателен');
      return;
    }

    setError(null);
    setSuccess(null);

    try {
      await dispatch(joinGroupByCode(groupCode)).unwrap();
      setSuccess('Вы успешно присоединились к группе!');
      setGroupCode('');
      closeForm();
    } catch (err: any) {
      setError(err || 'Не удалось присоединиться к группе. Попробуйте еще раз.');
    }
  };

  return (
    <div className='group-join-form'>
      <h2>Присоединиться к группе</h2>

      {error && <p style={{ color: 'red' }}>{error}</p>}
      {joinGroupByCodeError && <p style={{ color: 'red' }}>{joinGroupByCodeError}</p>}

      {success && <p style={{ color: 'green' }}>{success}</p>}
      <form onSubmit={handleFormSubmit}>
        <div>
          <label>Код группы:</label>
          <input
            type='text'
            value={groupCode}
            onChange={(e) => setGroupCode(e.target.value)}
            required
          />
        </div>
        <button type='submit' disabled={joinGroupByCodeLoading}>
          {joinGroupByCodeLoading ? 'Присоединение...' : 'Присоединиться'}
        </button>
      </form>
    </div>
  );
};

export default GroupJoinForm;
