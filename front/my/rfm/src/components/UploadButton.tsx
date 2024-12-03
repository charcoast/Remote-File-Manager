import * as React from 'react';
import { styled } from '@mui/material/styles';
import Button from '@mui/material/Button';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';

const VisuallyHiddenInput = styled('input')({
    clip: 'rect(0 0 0 0)',
    clipPath: 'inset(50%)',
    height: 1,
    overflow: 'hidden',
    position: 'absolute',
    bottom: 0,
    left: 0,
    whiteSpace: 'nowrap',
    width: 1,
});

interface UploadButtonProps {
    handleUpload: any
}

const UploadButtton: React.FC<UploadButtonProps> = ({ handleUpload }) => {
    return (
        <Button
            component="label"
            role={undefined}
            variant="contained"
            tabIndex={-1}
        >
            <CloudUploadIcon />
            <VisuallyHiddenInput
                type="file"
                onChange={handleUpload}
                multiple
            />
        </Button>
    );
}

export default UploadButtton
