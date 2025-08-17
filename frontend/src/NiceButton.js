
import { styled } from '@mui/material/styles';
import Button from '@mui/material/Button';


export const NiceButton = styled(Button)({
    width: "100%",
    height: "48px",
    boxShadow: 'none',
    textTransform: 'none',
    backgroundColor: "gray",
    padding: '1px',
    mt: "8px",
    lineHeight: 1,
    fontFamily: "monospace",
    '&:hover': {
        backgroundColor: '#0069d9',
        borderColor: '#0062cc',
        boxShadow: 'none',
    },
    '&:active': {
        boxShadow: 'none',
        backgroundColor: '#0062cc',
        borderColor: '#005cbf',
    },
    '&:focus': {
        boxShadow: '0 0 0 0.2rem rgba(0,123,255,.5)',
    },
    '&:disabled': {
        boxShadow: 'none',
        borderColor: '#005cbf',
        backgroundColor: 'gray',
    },
});

export const NiceVerticalButton = styled(Button)({
    width: "24px",
    height: "100vh",
    boxShadow: 'none',
    textTransform: 'none',
    backgroundColor: "gray",
    padding: '1px',
    mt: "8px",
    lineHeight: 1,
    fontFamily: "monospace",
    '&:hover': {
        backgroundColor: '#0069d9',
        borderColor: '#0062cc',
        boxShadow: 'none',
    },
    '&:active': {
        boxShadow: 'none',
        backgroundColor: '#0062cc',
        borderColor: '#005cbf',
    },
    '&:focus': {
        boxShadow: '0 0 0 0.2rem rgba(0,123,255,.5)',
    },
    '&:disabled': {
        boxShadow: 'none',
        borderColor: '#005cbf',
        backgroundColor: 'gray',
    },
});
