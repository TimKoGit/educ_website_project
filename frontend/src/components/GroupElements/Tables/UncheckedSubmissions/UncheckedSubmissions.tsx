import * as React from 'react';
import { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useParams, useNavigate } from 'react-router-dom';
import { RootState, AppDispatch } from '../../../../redux/store';
import { fetchUncheckedSubmissionsByGroupId, fetchUncheckedSubmissionsByContestId } from '../../../../redux/slices/submissionsSlice';
import './UncheckedSubmissions.css';


const UncheckedSubmissions: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { uncheckedSubmissions, fetchUncheckedSubmissionsByGroupIdLoading, fetchUncheckedSubmissionsByGroupIdError } = useSelector((state: RootState) => state.submissions);
  const { groupid, contestid } = useParams<{ groupid: string; contestid: string }>();
  const { user } = useSelector((state: RootState) => state.auth);
  const navigate = useNavigate();

  useEffect(() => {
    if (groupid && user && user.role === "teacher") {
      dispatch(fetchUncheckedSubmissionsByGroupId({ groupid }));
    }
  }, [dispatch, groupid, user]);

  useEffect(() => {
    if (contestid && user && user.role === "teacher") {
      dispatch(fetchUncheckedSubmissionsByContestId({ contestid }));
    }
  }, [dispatch, contestid, user]);

  const handleRowClick = (submissionId: string) => {
    navigate(`/submission/${submissionId}`);
  };

  if (!user || user.role === 'student')
    return;

  return (
    <div className='unchecked-submissions'>
      {fetchUncheckedSubmissionsByGroupIdLoading && <p>Загрузка посылок...</p>}
      {fetchUncheckedSubmissionsByGroupIdError && <p className='error'>{fetchUncheckedSubmissionsByGroupIdError}</p>}
      {(groupid || contestid) && (
      <div className='table-container unchecked-submissions-container'>
      <table>
        <thead>
          <tr>
            <th>Непроверенные посылки</th>
          </tr>
        </thead>
        <tbody>
          {uncheckedSubmissions.map((submission) => (
            <tr
              key={submission.id}
              onClick={() => handleRowClick(submission.id)}
              className='clickable-row'
            >
              <td><strong>{submission.studentName}</strong> - {submission.taskName}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>)}
    </div>
  );
};


export default UncheckedSubmissions;
