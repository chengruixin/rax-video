import React, { useState } from 'react'
import { Container } from '@material-ui/core'
import './index.css'
import VideoCard from './VideoCard'
import { useStyles } from './style'
import { useHistory, useLocation, useRouteMatch } from 'react-router'
import { fetchVideos } from '../../services'
import { useEffect } from 'react'

function VideosView() {
    const state  = useVideos()
    const classes = useStyles()

    console.log(state);

    return (
        <>
            <Container maxWidth="lg">
                <ul className={classes.ul}>
                    {state &&
                        state.map((item) => (
                            <VideoCard information={item} key={item._id} />
                        ))}
                </ul>
            </Container>
        </>
    )
}

function useVideos() {
    const { state } = useLocation()
    const history = useHistory()
    const { location } = history
    const params = new URLSearchParams(location.search)
    const searchVal = params.get('q')

    if (state) {
        console.log('return', state);
        return state
    }

    if (!searchVal || searchVal.length === 0) {
        return null
    }

    fetchVideos(searchVal).then(({ videoInfos }) => {
        history.replace("/videos?q=" + searchVal, videoInfos);
    }).catch(console.error)

    return null;
    
}
export default VideosView
