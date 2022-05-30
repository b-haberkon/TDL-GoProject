import React, { useEffect, useState } from 'react';
import { Navigate } from 'react-router-dom';

const Home = (props: { name: string }) => {

    return (
        <div>
            { props.name  ? 'Hi ' + props.name : 'You are not logged in'}
        </div>
    );
};

export default Home;