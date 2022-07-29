import * as React from "react";
import ArticleList from "../../Components/Article/ArticleList";
import {render, screen, waitFor} from "@testing-library/react";
import {MemoryRouter} from "react-router-dom";

// Arrange: Making a fake response
const articleList = [
    {
        "articleId": "1",
        "mainVersionId": "1",
        "title": "Web-Development research",
        "date_created": "2018-03-29",
        "owners": ["Jane Doe", "John Doe"],
    },
    {
        "articleId":"2",
        "mainVersionId": "1",
        "title":"Machine learning development",
        "date_created": "2015-10-26",
        "owners":["Jane Doe", "John Doe"],
    }
]

let mockedFetch = null;

describe('ArticleList', () => {
    // Arrange: Mock the fetch
    beforeEach(() => {
        const response = {
            status: 200,
            ok: true,
            json: jest.fn().mockResolvedValue(articleList) };
        mockedFetch = jest.fn().mockResolvedValue(response);
        global.fetch = mockedFetch;
    })

    // Test whether the first article is a part of the articleList
    test("test articles", async () => {
        // Act: Render the article list
        render(
            <MemoryRouter>
                <ArticleList />
            </MemoryRouter>
        );

        // Wait until it is finished loading
        const loadingSpinner = await screen.findByTestId("loadingSpinner")
        await waitFor(() => {
            expect(loadingSpinner).not.toBeInTheDocument()
        })

        // Assert that the fetch was used
        expect(mockedFetch).toHaveBeenCalledTimes(1)

        // Perform assertions for each article
        articleList.forEach((article) => {
            const i = article.articleId
            expect(screen.getByTestId("title" + i)).toHaveTextContent(article.title)
            article.owners.forEach((owner)=> {
                expect(screen.getByTestId("owners" + i)).toHaveTextContent(owner)
            })
        })
    })
})

let alert = {
    message: "message"
}
describe('ArticleList error', () => {
    // Arrange: Mock the fetch
    beforeEach(() => {
        const response = {
            status: 404,
            ok: false,
            json: jest.fn().mockResolvedValue(alert) };
        mockedFetch = jest.fn().mockResolvedValue(response);
        global.fetch = mockedFetch;
    })

    // Test whether the first article is a part of the articleList
    test("test articles", async () => {
        // Act: Render the article list
        render(
            <MemoryRouter>
                <ArticleList />
            </MemoryRouter>
        );

        // Wait until it is finished loading
        const loadingSpinner = await screen.findByTestId("loadingSpinner")
        await waitFor(() => {
            expect(loadingSpinner).not.toBeInTheDocument()
        })

        // Assert that the fetch was used
        expect(mockedFetch).toHaveBeenCalledTimes(1)

        // Perform assertions for error message
        expect(screen.getByRole("alert")).toHaveTextContent(
            "Error: Something went wrong when getting the articles. Error: " + alert.message
        )
    })
})