import * as React from "react"
import { Container, Header, Segment } from 'semantic-ui-react'

const style = {
    h1: {
        marginTop: '3em',
    },
    h2: {
        margin: '4em 0em 2em',
    },
    h3: {
        marginTop: '2em',
        padding: '2em 0em',
    },
    last: {
        marginBottom: '300px',
    },
}

const events = [
    {
        date: '1 Hour Ago',
        image: '/images/avatar/small/elliot.jpg',
        meta: '4 Likes',
        summary: 'Elliot Fu added you as a friend',
    },
    {
        date: '4 days ago',
        image: '/images/avatar/small/helen.jpg',
        meta: '1 Like',
        summary: 'Helen Troy added 2 new illustrations',
        extraImages: [
            '/images/wireframe/image.png',
            '/images/wireframe/image-text.png',
        ],
    },
    {
        date: '3 days ago',
        image: '/images/avatar/small/joe.jpg',
        meta: '8 Likes',
        summary: 'Joe Henderson posted on his page',
        extraText:
            "Ours is a life of constant reruns. We're always circling back to where we'd we started.",
    },
    {
        date: '4 days ago',
        image: '/images/avatar/small/justen.jpg',
        meta: '41 Likes',
        summary: 'Justen Kitsune added 2 new photos of you',
        extraText:
            'Look at these fun pics I found from a few years ago. Good times.',
        extraImages: [
            '/images/wireframe/image.png',
            '/images/wireframe/image-text.png',
        ],
    },
];

export class EventFeed extends React.Component {
    constructor(props){
        super(props)
    }

    render() {
        return (
            <div>
            <Header as='h3' textAlign='center' style={style.h3} content='Container' />
                <Container>
                    <Segment.Group>
                    <Segment>Content</Segment>
                    <Segment>Content</Segment>
                    <Segment>Content</Segment>
                    <Segment>Content</Segment>
                    </Segment.Group>
                </Container>
            </div>
        )
    }
}